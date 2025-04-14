

export function showToast(message) {
  Toastify({
    text: message,
    duration: 5000,
    gravity: "top",
    close: true,
    position: "center",
    style: {
      background: "var(--card-background)",
      color: "var(--text-color)",
      padding: "1rem 0rem",
      paddingLeft: "2rem",
      paddingRight: "1rem",
      borderRadius: "0.5rem",
      boxShadow: "0 0.5rem 1rem rgba(0, 0, 0, 0.1)",
      border: "3px solid var(--border-color)",
      display: "flex",
      gap: "4rem",
    }
  }).showToast();
}