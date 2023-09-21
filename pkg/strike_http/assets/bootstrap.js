import { Router } from "./router.js";
import React from "react";
import { hydrateRoot } from "react-dom/client";
import { jsx, jsxs } from "react/jsx-runtime";
import { ErrorBoundary } from "react-error-boundary";

React.startTransition(() => {
  hydrateRoot(document, jsx(Root, {}));
});

function Root() {
  return jsx(ErrorBoundary, {
    FallbackComponent: FallbackError,
    children: jsx(Router, {}),
  });
}

function FallbackError({ error }) {
  return jsxs("html", {
    children: [
      jsx("head", {
        children: jsx("title", { children: "Error" }),
      }),
      jsxs("body", {
        children: [
          jsx("h1", { children: "Fatal Error" }),
          jsx("pre", {
            style: { whiteSpace: "pre-wrap" },
            children: error.stack,
          }),
        ],
      }),
    ],
  });
}
