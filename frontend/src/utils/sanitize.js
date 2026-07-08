export function safeText(value) {
  return String(value ?? '').replace(/[\u0000-\u001f\u007f]/g, '');
}

