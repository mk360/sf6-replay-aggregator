const dotenv = require("dotenv");
dotenv.config({
    path: "../.env"
});

const TheFGCPlace = require("../channels/The FGC Place.json");
const SF6HighLevelReplays = require("../channels/SF6 High Level Replays.json");
const fs = require("fs");

const channelNames = {
    "The FGC Place": {
        data: TheFGCPlace,
        transformerFunc(name) {
            return `SF6 🔥 ${name}`;
        }
    },
    "SF6 High Level Replays": {
        data: SF6HighLevelReplays,
        transformerFunc(name) {
            return `${name} ▰ high level gameplay [SF6]`;
        }
    }
};

const urlQueryParams = new URLSearchParams();
urlQueryParams.append("part", "snippet,status");
urlQueryParams.append("maxResults", "50");
urlQueryParams.append("key", process.env.API_KEY);


function sendRequest(queryParams) {
    const url = `https://www.googleapis.com/youtube/v3/playlistItems?${queryParams.toString()}`;
    return fetch(url).then((resp) => resp.json());
}

(async function fct() {
    await Promise.all(["Ken", "Ryu", "Chun-Li", "Ed", "A.K.I.", "M. Bison", "Guile", "JP", "Juri", "Akuma", "Blanka", "Cammy", "Dhalsim", "Kimberly", "Lily", "Luke", "Manon", "Marisa", "Rashid", "Zangief", "E.Honda"].map((character) => {
        return generateCharacterPlaylist(character);
    }));
})();

async function generateCharacterPlaylist(characterName) {
    let resultsArray = [];

    for (let channelName in channelNames) {
        const transformedName = channelNames[channelName].transformerFunc(characterName);
        const playlistId = channelNames[channelName].data[transformedName];

        if (playlistId) {
            urlQueryParams.set("playlistId", playlistId);

            const data = await sendRequest(urlQueryParams);
            resultsArray = resultsArray.concat(data.items.filter((video) => video.status.privacyStatus !== "private").map((item) => ({
                title: item.snippet.title,
                id: item.snippet.resourceId.videoId,
                thumbnail: item.snippet.thumbnails.medium.url
            })));
        } else {
            console.log("cannot find playlist for character " + transformedName)
        }
    }

    fs.writeFileSync(characterName + ".json", JSON.stringify(resultsArray));
}
