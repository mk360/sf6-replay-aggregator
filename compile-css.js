const fs = require("fs");
let compiledCss = "";
const staticDir = fs.readdirSync("./static");
const staticCss = staticDir.filter((file) => file.endsWith(".css") && !file.includes("compiled"));

for (let i = 0; i < staticCss.length; i++) {
    const style = fs.readFileSync(`./static/${staticCss[i]}`, "utf-8");
    compiledCss += style + "\n";
}

fs.writeFileSync("./static/compiled.css", compiledCss);

