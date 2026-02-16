/**
 * Card Expansion Manager
 * Handles interactive card expansion with countdown timer and animations
 */

class CardManager {
  constructor() {
    this.cards = document.querySelectorAll(".card");
    this.expandDelay = 500; // milliseconds
    this.init();
  }

  /**
   * Initialize card event listeners and ring elements
   */
  init() {
    this.cards.forEach(card => {
      this.setupCardRing(card);
      this.attachEventListeners(card);
    });

    this.attachGlobalListeners();
    this.attachResizeListener();
  }

  /**
   * Setup SVG countdown ring for a card
   */
  setupCardRing(card) {
    if (card.querySelector(".card-ring")) {
      return;
    }

    const ring = document.createElement("div");
    ring.className = "card-ring";
    ring.innerHTML =
      '<svg viewBox="0 0 100 100" preserveAspectRatio="none">' +
      '<rect x="1" y="1" width="98" height="98" rx="6" ry="6" />' +
      '</svg>';
    card.appendChild(ring);

    card.style.setProperty("--countdown", `${this.expandDelay}ms`);
    this.updateRingLength(card);
  }

  /**
   * Attach hover and click event listeners to card
   */
  attachEventListeners(card) {
    card.addEventListener("mouseenter", () => this.handleCardHover(card));
    card.addEventListener("mouseleave", () => this.collapseCard(card));
  }

  /**
   * Handle card hover - start countdown to expand
   */
  handleCardHover(card) {
    // Clear any existing timer
    if (card._expandTimer) {
      clearTimeout(card._expandTimer);
    }

    // Update visual state
    card.classList.add("is-hovered", "is-countdown");
    card.classList.remove("is-armed");

    // Start countdown
    card._expandTimer = setTimeout(() => {
      this.expandCard(card);
    }, this.expandDelay);
  }

  /**
   * Expand card to show full details
   */
  expandCard(card) {
    if (card.classList.contains("is-expanded")) {
      return;
    }

    const grid = card.closest(".grid");
    if (!grid) {
      return;
    }

    // Get current and target dimensions
    const cardRect = card.getBoundingClientRect();
    const gridRect = grid.getBoundingClientRect();
    const currentTop = cardRect.top - gridRect.top;
    const currentLeft = cardRect.left - gridRect.left;

    // Calculate target dimensions
    const targetDimensions = this.calculateTargetDimensions(
      card,
      grid,
      cardRect,
      currentTop,
      currentLeft
    );

    // Create placeholder to maintain grid layout
    const placeholder = this.createPlaceholder(cardRect);
    card.parentNode.insertBefore(placeholder, card);
    card._placeholder = placeholder;

    // Prepare card for expansion
    grid.classList.add("is-expanded");
    card.classList.add("is-expanded");
    card.classList.remove("is-countdown", "is-armed");
    card.style.position = "absolute";
    card.style.top = `${currentTop}px`;
    card.style.left = `${currentLeft}px`;
    card.style.width = `${cardRect.width}px`;
    card.style.height = `${cardRect.height}px`;

    // Animate to target dimensions
    requestAnimationFrame(() => {
      card.style.top = `${targetDimensions.top}px`;
      card.style.left = `${targetDimensions.left}px`;
      card.style.width = `${targetDimensions.width}px`;
      card.style.height = `${targetDimensions.height}px`;
    });
  }

  /**
   * Calculate target dimensions for expanded card
   */
  calculateTargetDimensions(card, grid, cardRect, currentTop, currentLeft) {
    const maxWidth = grid.clientWidth;
    const maxHeight = grid.clientHeight;
    const minWidth = 520;
    const minHeight = 480;
    const contentPadding = 40;

    // Calculate based on content
    let width = Math.max(minWidth, card.scrollWidth + contentPadding);
    let height = Math.max(minHeight, card.scrollHeight + contentPadding);

    // Apply maximum constraints
    width = Math.min(width, Math.max(maxWidth * 0.95, minWidth + 200));
    height = Math.min(height, Math.max(maxHeight * 0.9, minHeight + 300));

    // Center within grid
    const cardCenterX = currentLeft + cardRect.width / 2;
    const cardCenterY = currentTop + cardRect.height / 2;
    const left = Math.max(0, Math.min(cardCenterX - width / 2, maxWidth - width));
    const top = Math.max(0, Math.min(cardCenterY - height / 2, maxHeight - height));

    return { width, height, left, top };
  }

  /**
   * Create placeholder element to preserve grid spacing
   */
  createPlaceholder(cardRect) {
    const placeholder = document.createElement("div");
    placeholder.className = "card-placeholder";
    placeholder.style.width = `${cardRect.width}px`;
    placeholder.style.height = "180px";
    return placeholder;
  }

  /**
   * Collapse expanded card back to normal size
   */
  collapseCard(card) {
    const grid = card.closest(".grid");

    // Clear any pending expand timer
    if (card._expandTimer) {
      clearTimeout(card._expandTimer);
      card._expandTimer = null;
    }

    // If not expanded, just clear hover state
    if (!card.classList.contains("is-expanded")) {
      card.classList.remove("is-hovered", "is-countdown", "is-armed");
      return;
    }

    // Get placeholder and grid data
    const placeholder = card._placeholder;
    if (!grid || !placeholder) {
      this.cleanupCardState(card, grid);
      return;
    }

    // Calculate original position
    const gridRect = grid.getBoundingClientRect();
    const placeholderRect = placeholder.getBoundingClientRect();
    const originalTop = placeholderRect.top - gridRect.top + grid.scrollTop;
    const originalLeft = placeholderRect.left - gridRect.left + grid.scrollLeft;

    // Animate back to original size
    card.addEventListener(
      "transitionend",
      () => this.cleanupCardState(card, grid),
      { once: true }
    );

    card.style.top = `${originalTop}px`;
    card.style.left = `${originalLeft}px`;
    card.style.width = `${placeholderRect.width}px`;
    card.style.height = `${placeholderRect.height}px`;
  }

  /**
   * Clean up card state after collapse
   */
  cleanupCardState(card, grid) {
    const placeholder = card._placeholder;
    if (placeholder) {
      placeholder.remove();
      card._placeholder = null;
    }

    // Remove inline styles
    card.style.removeProperty("top");
    card.style.removeProperty("left");
    card.style.removeProperty("width");
    card.style.removeProperty("height");
    card.style.removeProperty("position");

    // Update classes
    card.classList.remove("is-expanded", "is-hovered", "is-countdown", "is-armed");
    if (grid) {
      grid.classList.remove("is-expanded");
    }
  }

  /**
   * Update SVG ring length for smooth animation
   */
  updateRingLength(card) {
    const ring = card.querySelector(".card-ring");
    if (!ring) {
      return;
    }

    const svg = ring.querySelector("svg");
    const rect = ring.querySelector("rect");
    if (!svg || !rect) {
      return;
    }

    const ringRect = ring.getBoundingClientRect();
    const { width, height } = {
      width: Math.max(0, ringRect.width),
      height: Math.max(0, ringRect.height)
    };
    const radius = 6;

    // Update SVG attributes
    svg.setAttribute("viewBox", `0 0 ${width} ${height}`);
    rect.setAttribute("x", "1");
    rect.setAttribute("y", "1");
    rect.setAttribute("width", `${Math.max(0, width - 2)}`);
    rect.setAttribute("height", `${Math.max(0, height - 2)}`);
    rect.setAttribute("rx", `${radius}`);
    rect.setAttribute("ry", `${radius}`);

    // Update CSS custom property with calculated path length
    const length = rect.getTotalLength();
    card.style.setProperty("--ring-length", `${length}`);
  }

  /**
   * Handle clicks outside expanded card to close it
   */
  attachGlobalListeners() {
    document.addEventListener("click", event => {
      this.cards.forEach(card => {
        if (card.classList.contains("is-expanded")) {
          if (!card.contains(event.target)) {
            this.collapseCard(card);
          }
        }
      });
    });
  }

  /**
   * Update ring lengths on window resize
   */
  attachResizeListener() {
    window.addEventListener("resize", () => {
      this.cards.forEach(card => this.updateRingLength(card));
    });
  }
}

// Initialize when DOM is ready
document.addEventListener("DOMContentLoaded", () => {
  new CardManager();
});
