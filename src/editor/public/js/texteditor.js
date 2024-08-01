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

function EnableDragging(element, spacerElement) {

    element.addEventListener('mousedown', function(e) {
        let x = e.clientX - element.getBoundingClientRect().left;
        let y = e.clientY - element.getBoundingClientRect().top;

        element.style.zIndex = 1000;
        element.style.position = 'absolute';
        element.style.setProperty('opacity', 0.8, 'important');
        element.style.animationDelay = '0s';

        spacerElement.style.display = 'flex';

        move(e);

        function move(e) {
            element.style.left = `${e.clientX - x}px`;
            element.style.top = `${e.clientY - y}px`;

            let instructions = document.elementsFromPoint(e.clientX, e.clientY).filter(e => e.classList.contains('instruction'));

            if (instructions.length < 2) {
                element.style.backgroundColor = 'red';
            } else {
                element.style.backgroundColor = '';

                element.parentNode.insertBefore(spacerElement, instructions[1]);
            }

        }

        function up(e) {
            document.removeEventListener('mousemove', move);
            document.removeEventListener('mouseup', up);




            // Elements that its touching
            let elements = document.elementsFromPoint(e.clientX, e.clientY);
            let instructions = elements.filter(e => e.classList.contains('instruction'));
            element.style.position = 'relative';
            element.style.zIndex = 0;
            element.style.opacity = 1;

            element.style.left = '';
            element.style.top = '';

            if (!instructions[1]) {
                element.style.backgroundColor = '';
                return;
            }
            
            element.parentNode.insertBefore(element, instructions[1]);

            spacerElement.style.display = 'none';

        }

        document.addEventListener('mousemove', move);
        document.addEventListener('mouseup', up);
    })
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
        newElement('span', instruction, {innerText: text});
        newElement('span', instruction, {innerText: args.join(' ')});

        // cool animation uwu
        instruction.style.animationDelay = `${Math.random() * 0.5}s`;

        EnableDragging(instruction, this.spacerElement);

        return instruction;
    }

    initialize() {
        // Remove the version from the instructions
        this.instructions.shift();
        this.render();
    }

    render() {
        this.clearInstructions();

        this.spacerElement = newElement('div', this.HTMLParent, {className: 'spacer'});
        for (let instruction of this.instructions) {
            this.newInstruction(instruction);
        }

        
    }
}
