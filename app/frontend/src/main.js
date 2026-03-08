import {AppService} from "../bindings/github.com/typemd/typemd/app";

const listElement = document.getElementById("object-list");

async function loadObjects() {
    try {
        const objects = await AppService.ListObjects();

        if (!objects || objects.length === 0) {
            listElement.innerHTML = '<p class="empty">No objects found. Create some objects in your vault first.</p>';
            return;
        }

        // Group objects by type
        const groups = {};
        for (const obj of objects) {
            if (!groups[obj.type]) {
                groups[obj.type] = { emoji: obj.emoji, items: [] };
            }
            groups[obj.type].items.push(obj);
        }

        // Render groups
        let html = "";
        for (const [typeName, group] of Object.entries(groups).sort((a, b) => a[0].localeCompare(b[0]))) {
            const emoji = group.emoji ? `${group.emoji} ` : "";
            html += `<section class="type-group">`;
            html += `<h2>${emoji}${typeName} <span class="count">(${group.items.length})</span></h2>`;
            html += `<ul>`;
            for (const item of group.items) {
                html += `<li class="object-item" data-id="${item.id}">${item.displayName}</li>`;
            }
            html += `</ul>`;
            html += `</section>`;
        }
        listElement.innerHTML = html;

        // Add click handlers
        document.querySelectorAll(".object-item").forEach((el) => {
            el.addEventListener("click", () => showObject(el.dataset.id));
        });
    } catch (err) {
        listElement.innerHTML = `<p class="error">Error loading objects: ${err.message || err}</p>`;
    }
}

async function showObject(id) {
    try {
        const obj = await AppService.GetObject(id);
        const detail = document.getElementById("object-detail");
        if (detail) detail.remove();

        const section = document.createElement("aside");
        section.id = "object-detail";

        let propsHtml = "";
        if (obj.properties) {
            for (const [key, value] of Object.entries(obj.properties)) {
                propsHtml += `<tr><td class="prop-key">${key}</td><td class="prop-value">${value}</td></tr>`;
            }
        }

        section.innerHTML = `
            <div class="detail-header">
                <h2>${obj.displayName}</h2>
                <span class="detail-type">${obj.type}</span>
                <button class="close-btn" onclick="this.closest('#object-detail').remove()">✕</button>
            </div>
            ${propsHtml ? `<table class="props-table">${propsHtml}</table>` : ""}
            ${obj.body ? `<div class="body-content">${obj.body}</div>` : ""}
        `;
        document.getElementById("app").appendChild(section);
    } catch (err) {
        console.error("Error loading object:", err);
    }
}

loadObjects();
