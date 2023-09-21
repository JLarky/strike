import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";

export function Counter({ serverCounter }) {
  const [count, setCount] = React.useState(0);
  return jsx("button", {
    onClick: () => {
      setCount((x) => x + 1);
    },
    children: `Count: ${count} (${serverCounter})`,
  });
}
