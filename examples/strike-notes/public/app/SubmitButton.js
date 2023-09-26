import { createElement, useEffect, useRef } from "react";
import { experimental_useFormStatus as useFormStatus } from "react-dom";
import { jsx } from "react/jsx-runtime";

export function SubmitButton(props) {
  const { pending } = useFormStatus();
  const child = props.children[0];
  const ref = useRef();
  return createElement(child.type, {
    ...child.props,
    ref,
    children: pending ? "Loading..." : child.props.children,
  });
}
