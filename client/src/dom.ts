function isEventProperty(name) {
  return /^on/.test(name);
}

function extractEventName(name) {
  return name.slice(2).toLowerCase();
}

function setProperty(target: HTMLElement, name: string, value: any) {
  if (isEventProperty(name)) {
    return;
  }
  target.setAttribute(name, value)
}

function setProperties(target: HTMLElement, props: {[name: string]: any}) {
  for (let key in props) {
    target.setAttribute(key, props[key]);
  }
}

function addEventListeners(target: HTMLElement, props: {[name: string]: any}) {
  for (let name in props) {
    if (isEventProperty(name)) {
      target.addEventListener(
        extractEventName(name),
        props[name]
      );
    }
  }
}

function addChildren(target: HTMLElement, children: (HTMLElement|string)[]) {
  for (let child of children) {
    if (typeof child === 'string') {
      target.appendChild(document.createTextNode(child));
    } else {
      target.appendChild(child)
    }
  }
}

export function createElement(typ: string,
                              props: {[name: string]: any},
                              ...children: (HTMLElement|string)[]) : HTMLElement {

  let el = document.createElement(typ)
  setProperties(el, props);
  addEventListeners(el, props);
  addChildren(el, children);

  return el
}
