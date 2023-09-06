import { StrikeIsland } from "./islands.js";
import React from "https://esm.sh/react@canary?dev";
import { hydrateRoot } from "https://esm.sh/react-dom@canary/client?dev";
import { jsx, jsxs } from "https://esm.sh/react@canary/jsx-runtime?dev";

renderPage(JSON.parse(__rsc[0]));

export function renderPage(jsonData) {
  function jsonToJSX(x) {
    if (typeof x === "string") {
      return x;
    }

    const children = x.children?.map(jsonToJSX);
    const node = jsxs(x.tag_type, { ...x.props, children });

    if (x.tag_type === "strike-island") {
      const exportName = x.props["component-export"];
      return jsx(StrikeIsland, { exportName, serverProps: x, children: node });
    }

    return node;
  }
  const page = jsonToJSX(jsonData);

  const ClientRouter = ({ initialUrl }) => {
    console.log("initialUrl", initialUrl);
    return page;
    // return jsx("div", {
    //   children: jsx("a", { href: "/about", children: "About" }),
    // });
  };

  const root = hydrateRoot(document, page);

  //   root.render(jsx(ClientRouter, { initialUrl: "/" }));

  let currentPathname = window.location.pathname;

  async function navigate(pathname) {
    currentPathname = pathname;
    const clientJSX = await fetchClientJSX(pathname);
    if (pathname === currentPathname) {
      root.render(clientJSX);
    }
  }

  async function fetchClientJSX(pathname) {
    const response = await fetch(pathname, {
      headers: { RSC: "1" },
    });
    const clientJSXString = await response.text();
    console.log("clientJSXString", clientJSXString);
    const clientJSX = jsonToJSX(JSON.parse(clientJSXString));
    return clientJSX;
  }

  window.addEventListener(
    "click",
    (e) => {
      if (e.target.tagName !== "A") {
        return;
      }
      if (e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) {
        return;
      }
      const href = e.target.getAttribute("href");
      if (!href.startsWith("/")) {
        return;
      }
      e.preventDefault();
      window.history.pushState(null, null, href);
      navigate(href);
    },
    true
  );

  window.addEventListener("popstate", () => {
    navigate(window.location.pathname);
  });
  window.rscNav = (pathname) => {
    window.history.pushState(null, null, pathname);
    navigate(pathname);
  };
}
