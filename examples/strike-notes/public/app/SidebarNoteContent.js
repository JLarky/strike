// app/SidebarNoteContent.jsx
import {useState, useRef, useEffect, useTransition} from "react";
import {
jsx
} from "react/jsx-runtime";
function SidebarNoteContent({
  id,
  title,
  children,
  expandedChildren
}) {
  const [isPending, startTransition] = useTransition();
  const [isExpanded, setIsExpanded] = useState(false);
  const isActive = window.location.pathname === "/" + id;
  const itemRef = useRef(null);
  const prevTitleRef = useRef(title);
  useEffect(() => {
    if (title !== prevTitleRef.current) {
      prevTitleRef.current = title;
      itemRef.current.classList.add("flash");
    }
  }, [title]);
  return jsx("div", {
    ref: itemRef,
    onAnimationEnd: () => {
      itemRef.current.classList.remove("flash");
    },
    className: [
      "sidebar-note-list-item",
      isExpanded ? "note-expanded" : ""
    ].join(" "),
    children: [
      children,
      jsx("button", {
        className: "sidebar-note-open",
        style: {
          backgroundColor: isPending ? "var(--gray-80)" : isActive ? "var(--tertiary-blue)" : "",
          border: isActive ? "1px solid var(--primary-border)" : "1px solid transparent"
        },
        onClick: () => {
          startTransition(() => {
            const q = new URLSearchParams(window.location.search).get("q");
            __rscNav(`/${id}` + (q ? `?q=${encodeURIComponent(q)}` : ""));
          });
        },
        children: "Open note for preview"
      }),
      jsx("button", {
        className: "sidebar-note-toggle-expand",
        onClick: (e) => {
          e.stopPropagation();
          setIsExpanded(!isExpanded);
        },
        children: isExpanded ? jsx("img", {
          src: "/static/chevron-down.svg",
          width: "10px",
          height: "10px",
          alt: "Collapse"
        }) : jsx("img", {
          src: "/static/chevron-up.svg",
          width: "10px",
          height: "10px",
          alt: "Expand"
        })
      }),
      isExpanded && expandedChildren
    ]
  });
}
export {
  SidebarNoteContent as default
};
