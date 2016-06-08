declare namespace CodeMirror {
  function runMode(contents: string, mode: any, output: any);
}

declare module 'codemirror/addon/runmode/runmode' {
  export = CodeMirror;
}
