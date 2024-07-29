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

        this.initialize();
    }

    clearInstructions() {
        this.HTMLParent.innerHTML = '';
    }

    newInstruction(text, args = []) {
        let instruction = newElement('div', this.HTMLParent, {className: 'instruction'});
        let instructionText = newElement('span', instruction, {innerText: text});
        let instructionArguments = newElement('span', instruction, {innerText: args.join(' ')});

        return instruction;
    }

    initialize() {
        // Remove the version from the instructions
        this.instructions.shift();


        this.render();
    }

    render() {
        this.clearInstructions();

        for (let instruction of this.instructions) {
            this.newInstruction(instruction);
        }
    }
}
