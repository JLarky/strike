/// <reference path="node_modules/bun-types/types.d.ts" />

(async function () {
  // to trigger bun --watch
  import("./app/EditButton.jsx");
  import("./app/SearchField.jsx");
  import("./app/SidebarNoteContent.jsx");
})();

const x = await Bun.build({
  entrypoints: [
    "app/EditButton.jsx",
    "app/SearchField.jsx",
    "app/SidebarNoteContent.jsx",
  ],
  external: ["react", "react-dom", "./framework/router.js"],
  outdir: "app",
});

if (!x.success) {
  console.log(x);
} else {
  console.log("Build successful");
}
