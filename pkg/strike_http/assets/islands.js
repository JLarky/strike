import * as islands from "strike_islands";
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";

export function StrikeIsland(props) {
  const { "component-export": exportName, ssrFallback, ...rest } = props;
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
    return jsx(comp, rest);
  }
  return ssrFallback;
}
