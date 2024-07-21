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


function sendRequest(queryParams) {
    const url = `https://www.googleapis.com/youtube/v3/playlistItems?${queryParams.toString()}`;
    return fetch(url).then((resp) => resp.json());
}

(async function fct() {
    console.time("js")
    await Promise.all(["Ken", "Ryu", "Chun-Li", "Ed", "A.K.I.", "Dee Jay", "M. Bison", "Guile", "JP", "Juri", "Akuma", "Blanka", "Cammy", "Dhalsim", "Kimberly", "Lily", "Luke", "Manon", "Marisa", "Rashid", "Zangief", "E.Honda"].map((character) => {
        return generateCharacterPlaylist(character);
    }));
    console.timeEnd("js")
})();

async function generateCharacterPlaylist(characterName) {
    for (let channelName in channelNames) {
        const transformedName = channelNames[channelName].transformerFunc(characterName);
        const playlistId = channelNames[channelName].data[transformedName];

        if (playlistId) {
            const urlQueryParams = new URLSearchParams();
            urlQueryParams.append("part", "snippet,status");
            urlQueryParams.append("maxResults", "50");
            urlQueryParams.append("key", process.env.API_KEY);
            urlQueryParams.set("playlistId", playlistId);
            const resultsArray = await fetchVideosIterative(urlQueryParams);
            fs.writeFileSync(characterName + ".json", JSON.stringify(resultsArray));
        } else {
            console.log("cannot find playlist for character " + transformedName)
        }
    }
}

async function fetchVideosIterative(searchParams) {
    let currentVideos = [];
    let nextPageToken = "";
    do {
        const data = await sendRequest(searchParams);
        currentVideos = currentVideos.concat(data.items.filter((video) => video.status.privacyStatus === "public").map(
            (item) => ({
            title: item.snippet.title,
            id: item.snippet.resourceId.videoId,
            thumbnail: item.snippet.thumbnails.medium.url
        })));
        nextPageToken = data.nextPageToken;
        if (nextPageToken) {
            searchParams.set("pageToken", data.nextPageToken);
        } else break;
    } while (nextPageToken)

    return currentVideos;
}
