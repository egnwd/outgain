export function createElement(typ: string, props: {[key: string]: string}, ...children) {
  let el = document.createElement(typ)

  for (let child of children) {
    if (typeof child === 'string') {
      child = document.createTextNode(child);
    }
    el.appendChild(child)
  }

  for (let key in props) {
    el.setAttribute(key, props[key]);
  }

  return el
}
