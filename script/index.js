(function() {
    const characterSelect = document.getElementById("character-select-list");

    let p1 = "ryu";
    let p2 = "ken";
    let chosenSide = "p1";

    const characterButtons = document.querySelectorAll(".character-grid button");

    const p1CharacterSelect = document.getElementById("p1");
    const p2CharacterSelect = document.getElementById("p2");

    p1CharacterSelect.onclick = function() {
        characterSelect.classList.remove("p2-select");
        if (chosenSide !== "p1") {
            if (!characterSelect.classList.contains("show"))  {
                characterSelect.classList.add("show");
            }
        } else {
            toggleClass(characterSelect, "show");
        }
        chosenSide = "p1";
    };

    p2CharacterSelect.onclick = function() {
        if (chosenSide !== "p2") {
            if (!characterSelect.classList.contains("show"))  {
                characterSelect.classList.add("show");
            }
            characterSelect.classList.add("p2-select");
        } else {
            toggleClass(characterSelect, "show");
        }
        chosenSide = "p2";
    };

    for (let i = 0; i < characterButtons.length; i++) {
        const button = characterButtons.item(i);
        button.onclick = function() {
            const characterName = this.id;
            updateCurrentCharacter(chosenSide, characterName);
            toggleClass(characterSelect, "show");
        }
    }

    function loadVideos() {
        
    }

    function updateCurrentCharacter(side, character) {
        const targetElement = side === "p1" ? p1CharacterSelect : p2CharacterSelect;
        const nameElement = targetElement.getElementsByClassName("name")[0];
        nameElement.textContent = character.replace("_", ". ");

        const characterImage = targetElement.getElementsByClassName("character-image")[0];
        const oldCharacterClass = Array.prototype.find.call(characterImage.classList, (l) => l.includes("-thumbnail"));
        characterImage.classList.remove(oldCharacterClass);
        characterImage.classList.add(`${character}-thumbnail`);
    }

    function toggleClass(element, cssClass) {
        if (Array.prototype.includes.call(element.classList, cssClass)) {
            element.classList.remove(cssClass);
        } else {
            element.classList.add(cssClass);
        }
    }
})();