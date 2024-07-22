const fs = require("fs");
const dotenv = require("dotenv");
dotenv.config({ path: "./server/.env" });

fs.writeFileSync("./script/api.js", `const API_URL="${process.env.API_URL}";`);
