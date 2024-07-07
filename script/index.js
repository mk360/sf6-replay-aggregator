(function() {
    const characterSelect = document.getElementById("character-select-list");

    let p1 = "ryu";
    let p2 = "ken";
    let chosenSide = "p1";

    const characterButtons = document.querySelectorAll(".character-grid button");
    const videosSection = document.getElementById("videos");
    const currentVideoContainer = document.getElementById("current-video");
    const p1CharacterSelect = document.getElementById("p1");
    const p2CharacterSelect = document.getElementById("p2");
    let currentVideoIframe = null;

    const specialFormatNames = {
        "ehonda": "E.Honda",
        "jp": "JP",
        "m-bison": "M. Bison",
        "dee-jay": "Dee Jay",
        "chun-li": "Chun-Li",
        "aki": "A.K.I."
    };

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

            if (chosenSide === "p1") {
                p1 = characterName;
            } else {
                p2 = characterName;
            }

            loadVideos();
        }
    }

    function formatName(name) {
        const formattedName = name in specialFormatNames ? specialFormatNames[name] : name[0].toUpperCase() + name.substring(1, name.length);
        return encodeURIComponent(formattedName);
    }

    function loadVideos() {
        fetch(`http://localhost:4444/replays?character=${formatName(p1)}`).then((resp) => resp.text()).then((responseText) => {
            videosSection.innerHTML = responseText;
            const videoPreviews = videosSection.querySelectorAll(".video-preview-bg");
            for (let i = 0; i < videoPreviews.length; i++) {
                const preview = videoPreviews[i];
                const youtubeId = preview.attributes["data-youtube-id"].value;
                preview.onclick = function() {
                    if (!currentVideoIframe) {
                        currentVideoIframe = document.createElement("iframe");
                        currentVideoContainer.appendChild(currentVideoIframe);
                        currentVideoIframe.classList.add("main-video");
                        currentVideoIframe.frameborder="0"
                    }
                    currentVideoIframe.src = `https://youtube.com/embed/${youtubeId}`;
                };
            }
            videosSection.classList.add("show-results");
        });
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