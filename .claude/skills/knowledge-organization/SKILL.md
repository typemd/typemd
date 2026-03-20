---
name: knowledge-organization
description: Use when designing knowledge organization features — type schemas, relations, labeling, hierarchy, taxonomy, tagging, metadata, or cross-vault mapping. Use when the user mentions "SKOS", "Dublin Core", "Schema.org", "taxonomy", "concept hierarchy", "broader/narrower", "preferred label", "alias", "metadata standard", "knowledge organization", or when evaluating whether a data model aligns with established knowledge management standards.
---

# Knowledge Organization Standards Reference

Reference for three standards relevant to typemd's data model design: **SKOS** (concept relationships), **Dublin Core** (metadata), and **Schema.org** (structured data). Each covers a different concern; together they provide a comprehensive vocabulary for knowledge organization.

## SKOS — Simple Knowledge Organization System

W3C standard for representing thesauri, taxonomies, and subject heading systems. Focuses on **concepts and their relationships**.

### Core Concepts

| Concept | Description |
|---------|-------------|
| `skos:Concept` | Basic unit of knowledge (an idea, a term) |
| `skos:ConceptScheme` | A collection of related concepts (e.g. "Programming Languages") |
| `skos:prefLabel` | Preferred label — one per language per concept |
| `skos:altLabel` | Alternative label (aliases, abbreviations, translations) |
| `skos:hiddenLabel` | Hidden label (common misspellings — aids search, not displayed) |
| `skos:definition` | Formal definition of the concept |
| `skos:inScheme` | Which concept scheme a concept belongs to |
| `skos:broader` | Broader concept ("Go" → "Programming Language") |
| `skos:narrower` | Narrower concept ("Programming Language" → "Go") |
| `skos:related` | Associative relation ("Order" ↔ "Payment") |

### Documentation Properties

| Property | Purpose |
|----------|---------|
| `skos:scopeNote` | When to use / not use this term |
| `skos:example` | Concrete usage examples |
| `skos:historyNote` | How the term evolved over time |
| `skos:changeNote` | Change log entries |
| `skos:editorialNote` | Editorial remarks (pending items) |

### Mapping Properties (cross-vocabulary)

| Property | Purpose |
|----------|---------|
| `skos:exactMatch` | Exact equivalence across vocabularies |
| `skos:closeMatch` | Near equivalence |
| `skos:broadMatch` / `narrowMatch` | Cross-vocabulary broader/narrower |
| `skos:relatedMatch` | Cross-vocabulary association |

### References

- [SKOS Reference (W3C Recommendation)](https://www.w3.org/TR/skos-reference/) — Normative specification
- [SKOS Primer (W3C Working Group Note)](https://www.w3.org/TR/skos-primer/) — Companion user guide (examples are illustrative, not normative)

## Dublin Core — Metadata Standard

DCMI standard for describing resources with structured metadata. Focuses on **who, what, when, where** of a resource. SKOS explicitly defers metadata concerns to Dublin Core.

### Core 15 Elements

| Element | Description |
|---------|-------------|
| `dc:title` | Name given to the resource |
| `dc:creator` | Entity primarily responsible for making the resource |
| `dc:subject` | Topic of the resource (controlled vocabulary recommended) |
| `dc:description` | Account of the resource (abstract, summary) |
| `dc:publisher` | Entity making the resource available |
| `dc:contributor` | Entity responsible for contributions |
| `dc:date` | Point or period associated with lifecycle events |
| `dc:type` | Nature or genre of the resource |
| `dc:format` | File format, physical medium, or dimensions |
| `dc:identifier` | Unambiguous reference (ISBN, DOI, URN) |
| `dc:source` | Related resource from which this derives |
| `dc:language` | Language of the resource |
| `dc:relation` | A related resource |
| `dc:coverage` | Spatial or temporal topic |
| `dc:rights` | Intellectual property information |

### Key Extended Terms (`dct:` namespace)

| Term | Description |
|------|-------------|
| `dct:created` | Date of creation |
| `dct:modified` | Date of last change |
| `dct:issued` | Date of formal publication |
| `dct:isPartOf` | Resource is part of another resource |
| `dct:hasPart` | Resource includes another resource |
| `dct:license` | Legal document granting permission |
| `dct:provenance` | Changes in ownership/custody |
| `dct:accessRights` | Who may access the resource |

### References

- [DCMI Metadata Terms](https://www.dublincore.org/specifications/dublin-core/dcmi-terms/)

## Schema.org — Structured Data Vocabulary

Collaborative vocabulary for structured data on the web. Focuses on **entity types and their properties** for interoperability.

### Key Types

| Type | Description |
|------|-------------|
| `Thing` | Base type for all entities (`name`, `description`, `identifier`, `alternateName`, `sameAs`, `url`) |
| `CreativeWork` | Content with authorship and lifecycle (extends Thing) |
| `DefinedTerm` | A word, name, or phrase with a formal definition |
| `DefinedTermSet` | Collection container for DefinedTerms (glossary, taxonomy) |
| `CategoryCode` | Subtype of DefinedTerm for hierarchical classification |
| `Person` / `Organization` | Entity types for authorship and attribution |

### CreativeWork Properties (knowledge-relevant)

| Property | Type | Description |
|----------|------|-------------|
| `author` | Person / Organization | Creator of the content |
| `dateCreated` | Date | When first produced |
| `dateModified` | Date | When last changed |
| `keywords` | DefinedTerm / Text | Tags describing the item |
| `about` | Thing | Subject matter |
| `isPartOf` / `hasPart` | CreativeWork | Part-whole relationships |
| `isBasedOn` | CreativeWork | Source material |
| `citation` | CreativeWork / Text | References to other works |
| `mentions` | Thing | Entities mentioned in the content |
| `creativeWorkStatus` | DefinedTerm / Text | Lifecycle stage (draft, published) |

### DefinedTerm Properties

| Property | Type | Description |
|----------|------|-------------|
| `name` | Text | The term being defined |
| `description` | Text | Definition or explanation |
| `alternateName` | Text | Alternative label |
| `inDefinedTermSet` | DefinedTermSet | Container set for this term |
| `termCode` | Text | Code identifying the term within its set |
| `sameAs` | URL | Link to identical entity elsewhere |

### References

- [Schema.org Full Hierarchy](https://schema.org/docs/full.html)
- [DefinedTerm](https://schema.org/DefinedTerm)
- [CreativeWork](https://schema.org/CreativeWork)

## Mapping to typemd

| Concern | SKOS | Dublin Core | Schema.org | typemd |
|---------|------|-------------|------------|--------|
| Name | `prefLabel` | `dc:title` | `name` | `name` |
| Definition | `definition` | `dc:description` | `description` | `description` |
| Alternative name | `altLabel` | — | `alternateName` | — |
| Identifier | — | `dc:identifier` | `identifier` | ObjectID (ULID) |
| Type/Category | `inScheme` | `dc:type` | `@type` | Object's `Type` field |
| Tags/Subject | — | `dc:subject` | `keywords` | `tags` |
| Created | — | `dct:created` | `dateCreated` | `created_at` |
| Modified | — | `dct:modified` | `dateModified` | `updated_at` |
| Association | `related` | `dc:relation` | — | `relation` property |
| Hierarchy | `broader`/`narrower` | — | `CategoryCode` | — |
| Part-whole | — | `isPartOf`/`hasPart` | `isPartOf`/`hasPart` | — |
| Source/Origin | — | `dc:source` | `isBasedOn` | — |
| Mentions | — | — | `mentions` | Wiki-links |
| Cross-system equivalence | `exactMatch`/`closeMatch` | — | `sameAs` | — |

## Concepts Not Applicable to typemd

| Concept | Reason |
|---------|--------|
| RDF / Turtle / JSON-LD syntax | typemd uses Markdown + YAML |
| Language tags (`@zh-TW`) | Personal/small-team tool, no i18n needed |
| `skos:notation` | ULID already solves unique identification |
| `dc:format` / `dc:language` / `dc:rights` | All objects are Markdown, single language, local-first |
| `dc:publisher` / `dc:coverage` | Not relevant for personal knowledge base |
| `hasTopConcept` | View's `group_by` provides similar capability |
