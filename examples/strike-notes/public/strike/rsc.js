// @ts-check
import React from "https://esm.sh/react@canary?dev";
import { jsx, jsxs } from "https://esm.sh/react@canary/jsx-runtime?dev";

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
  const clientJSX = jsonToJSX(JSON.parse(clientJSXString));
  console.log("clientJSX", clientJSX);
  return clientJSX;
});

export function jsonToJSX(x) {
  if (typeof x === "string") {
    return x;
  }
  // meta tags has to have undefined children instead of empty array
  let children = undefined;
  if (x?.children && x.children.length > 0) {
    children = x.children.map(jsonToJSX);
  }
  const props = {};
  props.children = children;
  for (const [k, v] of Object.entries(x.props || {})) {
    if (x.tag_type === "meta" && k === "charset") {
      props.charSet = v;
    } else {
      props[k] = v;
    }
  }
  const node = jsxs(x.tag_type, props);

  if (x.tag_type === "strike-island") {
    return jsx(StrikeIsland, { children: node });
  }

  return node;
}
