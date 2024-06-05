function newElement(tag, parent = null, options = {}) {
    let element = document.createElement(tag);
    for (let key in options) {
        element[key] = options[key];
    }

    if (parent) {
        parent.appendChild(element);
    }

    return element;
}

class Editor {
    constructor(HTMLParent, instructions = []) {
        this.instructions = instructions;
        this.HTMLParent = HTMLParent;
    }
}
