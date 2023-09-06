import * as islands from "./islands.js";
import React from "https://esm.sh/react@canary?dev";
import { hydrateRoot } from "https://esm.sh/react-dom@canary/client?dev";
import { jsx, jsxs } from "https://esm.sh/react@canary/jsx-runtime?dev";

console.log("counter");

export function StrikeIsland({ exportName, children }) {
  if (!exportName) {
    throw new Error(`strike-island is missing component-export prop`);
  }
  const comp = islands[exportName];
  if (!comp) {
    throw new Error(`island ${exportName} doesn't exist`);
  }
  const [isMounted, setIsMounted] = React.useState(false);
  React.useEffect(() => {
    setIsMounted(true);
  }, []);
  if (isMounted) {
    return jsx(comp, {});
  }
  return children;
}

export function Counter() {
  const [count, setCount] = React.useState(0);
  return jsx("button", {
    onClick: () => {
      setCount((x) => x + 1);
    },
    children: "Count: " + count,
  });
}
