import * as islands from "./islands.js";
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";

export { default as EditButton } from "../app/EditButton.js";
export { default as SearchField } from "../app/SearchField.js";
export { default as SidebarNoteContent } from "../app/SidebarNoteContent.js";

export function StrikeIsland(props) {
  const { "component-export": exportName } = props;
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
    return jsx("strike-island", { ...props, children: jsx(comp, props) });
  }
  return jsx("strike-island", props);
}
