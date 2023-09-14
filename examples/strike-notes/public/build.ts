/// <reference path="node_modules/bun-types/types.d.ts" />

const x = await Bun.build({
  entrypoints: ["app/EditButton.jsx", "app/SearchField.jsx"],
  external: ["react", "react-dom", "./framework/router.js"],
  outdir: "app",
});

if (!x.success) {
  console.log(x);
} else {
  console.log("Build successful");
}
