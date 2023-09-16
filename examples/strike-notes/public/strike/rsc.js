// @ts-check
import React from "react";
import { jsx, jsxs } from "react/jsx-runtime";
import { StrikeIsland } from "./islands.js";
import { StrikeSuspense } from "./suspense.js";

export function RscComponent({ initialPage, url, routerKey }) {
  if (initialPage) {
    return initialPage;
  }
  // return React.use(fetchClientJSX(url, routerKey));
  // Error: Support for `use` not yet implemented in react-debug-tools.
  return fetchClientJSX(url, routerKey);
}

const fetchClientJSX = React.cache(async function fetchClientJSX(
  pathname,
  key
) {
  const response = await fetch(pathname, {
    headers: { RSC: "1" },
  });
  const clientJSXString = await response.text();
  const clientJSX = jsonToJSX(clientJSXString);
  console.log("clientJSX", clientJSX);
  return clientJSX;
});

export function jsonToJSX(x) {
  const ctx = {};
  const parsed = JSON.parse(x, function fromJSON(key, value) {
    return parseModelString(ctx, this, key, value);
  });
  // console.log("str", parsed);
  return parsed;
}

/** @type {import("./rsc").parseModelString} */
function parseModelString(ctx, parent, key, value) {
  if (key === "$strike" && value === "component") {
    parent["$$typeof"] = Symbol.for("react.element");
    parent.type = parent["$type"];
    delete parent["$type"];
    parent.ref = null;
    parent.key = null;
    for (const [k, v] of Object.entries(parent.props || {})) {
      if (k === "style" && typeof v === "string") {
        /** @type {{ [key: string]: string }} */
        const style = {};
        v.split(";").forEach((x) => {
          const [k, v] = x.split(":");
          if (k && v) {
            style[k.trim()] = v.trim();
          }
        });
        parent.props.style = style;
      } else if (k === "key") {
        key = v;
        delete parent.props.key;
      } else if (k === "class") {
        delete parent.props.class;
        parent.props.className = v;
      } else if (parent.type === "meta" && k === "charset") {
        delete parent.props.charset;
        parent.props.charSet = v;
      } else {
        parent.props[k] = v;
      }
    }
    if (parent.type === "strike-suspense") {
      parent.type = StrikeSuspense;
    }
    if (parent.type === "strike-island") {
      parent.type = StrikeIsland;
    }
    return undefined;
  }
  return value;
}
