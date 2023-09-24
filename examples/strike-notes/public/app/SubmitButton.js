import { createElement, useState, useTransition } from "react";
import { experimental_useFormStatus as useFormStatus } from "react-dom";
import { jsx } from "react/jsx-runtime";

export function SubmitButton(props) {
  const { pending } = useFormStatus();
  console.log(pending);
  const child = props.children[0];
  return createElement(child.type, child.props);
}
