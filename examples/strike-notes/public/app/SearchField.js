// app/SearchField.jsx
import {useState, useTransition} from "react";

// app/Spinner.jsx
import {
jsx
} from "react/jsx-runtime";
function Spinner({ active = true }) {
  return jsx("div", {
    className: ["spinner", active && "spinner--active"].join(" "),
    role: "progressbar",
    "aria-busy": active ? "true" : "false"
  });
}

// app/SearchField.jsx
import {
jsx as jsx2
} from "react/jsx-runtime";
function SearchField() {
  const [text, setText] = useState(() => {
    const q = new URLSearchParams(window.location.search).get("q");
    return q || "";
  });
  const [isSearching, startSearching] = useTransition();
  return jsx2("form", {
    className: "search",
    role: "search",
    action: (x) => console.log(x),
    children: [
      jsx2("label", {
        className: "offscreen",
        htmlFor: "sidebar-search-input",
        children: "Search for a note by title"
      }),
      jsx2("input", {
        id: "sidebar-search-input",
        placeholder: "Search",
        value: text,
        onChange: (e) => {
          const newText = e.target.value;
          setText(newText);
          startSearching(() => {
            const url = new URL(window.location);
            url.searchParams.set("q", newText);
            __rscNav(url.toString());
          });
        }
      }),
      jsx2(Spinner, {
        active: isSearching
      })
    ]
  });
}
export {
  SearchField as default
};
