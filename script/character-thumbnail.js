const template = document.createElement("template");

class CharacterThumbnail extends HTMLElement {
  characterImage;
  static get observedAttributes() {
    return ["character"];
  }

  updateCharacterClass(newClass) {
    this.characterImage.classList.remove(
      Array.prototype.find.call(this.characterImage.classList, (cl) =>
        cl.includes("thumbnail")
      )
    );
    this.characterImage.classList.add(newClass);
  }

  attributeChangedCallback(attr, oldVal, newVal) {
    if (oldVal === newVal || !newVal.endsWith("-thumbnail")) return;
    this.updateCharacterClass(newVal);
  }

  constructor() {
    super();
    this._shadowRoot = this.attachShadow({ mode: "open" });
    const shadowDiv = document.createElement("div");
    shadowDiv.classList.add("shadow-layer");
    shadowDiv.classList.add("p1");
    const thumbnailContainer = document.createElement("div");
    thumbnailContainer.classList.add("thumbnail");
    const backgroundSvg = document.createElement("svg");
    backgroundSvg.viewbox = "0 0 660 720";
    for (let i = 0; i < 30; i++) {
      const rect = document.createElementNS(
        "http://www.w3.org/2000/svg",
        "rect"
      );
      rect.setAttribute("height", "20");
      rect.setAttribute("width", "800");
      rect.setAttribute("x", "40");
      rect.setAttribute("y", i * 25);
      backgroundSvg.appendChild(rect);
    }
    thumbnailContainer.appendChild(backgroundSvg);
    const actualImage = document.createElement("img");
    actualImage.id = "character-img";
    actualImage.classList.add("character-image");
    this.characterImage = actualImage;
    shadowDiv.appendChild(thumbnailContainer);
    this._shadowRoot.appendChild(shadowDiv);
  }

  connectedCallback() {
    this.updateCharacterClass(this.getAttribute("character"));
    if (this.getAttribute("side") === "p2") {
      this.characterImage.classList.add("p2-thumbnail");
    }
  }
}

customElements.define("character-thumbnail", CharacterThumbnail);
