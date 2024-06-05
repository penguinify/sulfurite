class Editor {
    constructor(HTMLParent, instructions = []) {
        this.instructions = instructions;
        this.HTMLParent = HTMLParent;
    }

    Initialize() {
        this.Container = document.createElement("div");
        this.Container.classList.add("editor");
        this.Container.classList.add("full");
        this.HTMLParent.appendChild(this.Container);
    }

    Render() {
        for (let i = 0; i < this.instructions.length; i++) {
            let line = document.createElement("div");
            line.classList.add("line");
            line.innerText = this.instructions[i];
            line.setAttribute("contenteditable", "true");
            this.Container.appendChild(line);
        }
    }
}
