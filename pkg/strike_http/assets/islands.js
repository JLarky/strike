import * as islands from "strike_islands";
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";

export function StrikeIsland(props) {
  const { component: exportName, ssrFallback, islandProps } = props;
  console.log("StrikeIsland", props);
  if (!exportName) {
    throw new Error(`strike-island is missing component prop`);
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
    return jsx(comp, islandProps || {});
  }
  return ssrFallback;
}
