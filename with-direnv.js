// A script to run the passed command, but calling direnv allow first
// This is needed because lint-staged runs in a subshell, and direnv
// doesn't work in subshells.

const { execSync } = require("child_process");

const command = process.argv.slice(2).join(" ");

const script = `
direnv allow;
eval "$(direnv export bash)";
${command};
`.trim();

execSync(script, { stdio: "inherit", env: process.env });
