package core

import "fmt"

// reconcileWikiLinksAndTags reconciles wikilinks and tag relations, emitting events.
func (r *Reconciler) reconcileWikiLinksAndTags(ctx *syncContext, result *ReconcileResult) ([]DomainEvent, error) {
	var events []DomainEvent
	modifiedBodies := make(map[string]string)

	for id, body := range ctx.diskBodies {
		obj := ctx.diskObjects[id]
		sourceType := ""
		if obj != nil {
			sourceType = obj.Type
		}
		wikiEvents, err := r.reconcileWikiLinks(id, body, sourceType, ctx, result, modifiedBodies)
		if err != nil {
			return nil, fmt.Errorf("reconcile wikilinks for %s: %w", id, err)
		}
		events = append(events, wikiEvents...)
	}

	// Write back expanded bodies to disk
	for id, newBody := range modifiedBodies {
		obj := ctx.diskObjects[id]
		if obj == nil {
			continue
		}
		obj.Body = newBody
		schema := ctx.schemas[obj.Type]
		keyOrder := OrderedPropKeys(obj.Properties, schema)
		if err := r.repo.Save(obj, keyOrder); err != nil {
			return nil, fmt.Errorf("write-back expanded wiki-links for %s: %w", id, err)
		}
	}

	tagEvents, err := r.reconcileTagRelations(ctx)
	if err != nil {
		return nil, err
	}
	events = append(events, tagEvents...)
	return events, nil
}

// reconcileWikiLinks extracts wiki-links from body, resolves shorthand targets,
// and emits WikiLinksSynced events.
func (r *Reconciler) reconcileWikiLinks(objectID, body, sourceType string, ctx *syncContext, result *ReconcileResult, modifiedBodies map[string]string) ([]DomainEvent, error) {
	links := ParseWikiLinks(body)
	if len(links) == 0 {
		return []DomainEvent{WikiLinksSynced{ObjectID: objectID, Links: nil}}, nil
	}

	entries := make([]WikiLinkEntry, len(links))
	resolutions := make(map[string]string)
	for i, link := range links {
		res := resolveWikiLinkTarget(link.Target, sourceType, ctx.diskIDs, ctx.nameIndex)
		entries[i] = WikiLinkEntry{
			ToID:        res.resolvedID,
			Target:      link.Target,
			DisplayText: link.DisplayText,
		}
		if res.changed && res.resolvedID != "" {
			resolutions[link.Target] = res.resolvedID
			result.WikiLinksExpanded++
		}
		if res.err != nil {
			result.UnresolvedWikiLinks = append(result.UnresolvedWikiLinks, UnresolvedWikiLink{
				ObjectID: objectID,
				Target:   link.Target,
				Reason:   resolveErrorReason(res.err),
				Matches:  resolveErrorMatches(res.err),
			})
		}
	}

	if len(resolutions) > 0 {
		newBody, _ := expandWikiLinksInBody(body, resolutions)
		modifiedBodies[objectID] = newBody
	}

	return []DomainEvent{WikiLinksSynced{ObjectID: objectID, Links: entries}}, nil
}
