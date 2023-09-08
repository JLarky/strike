// @ts-check
import { Router } from "./router.js";
import React from "https://esm.sh/react@canary?dev";
import { hydrateRoot } from "https://esm.sh/react-dom@canary/client?dev";
import { jsx, jsxs } from "https://esm.sh/react@canary/jsx-runtime?dev";
import { ErrorBoundary } from "https://esm.sh/react-error-boundary";

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
