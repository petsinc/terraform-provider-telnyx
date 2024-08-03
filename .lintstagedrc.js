module.exports = {
  "*.{tf,hcl}": "node with-direnv.js just format-hcl",
  "*.go": "node with-direnv.js just format-go",
};
