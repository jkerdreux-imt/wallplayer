// Player controls
function updatePlayPauseButton(paused) {
  const btn = document.getElementById("playPause");
  if (btn) {
    btn.querySelector(".material-symbols-outlined").textContent = paused ? "play_arrow" : "pause";
  }
}

function togglePlayPause() {
  const video = document.querySelector("video");
  if (!video) return;

  if (video.paused) {
    video.play();
  } else {
    video.pause();
  }
  // Button will be updated by play/pause event listeners
}

function toggleVolumeMenu(event) {
  if (event) {
    event.stopPropagation();
  }
  const control = document.querySelector(".volume-control");
  control.classList.toggle("open");
}

function setVolume(value) {
  // Store in localStorage
  localStorage.setItem("volume", value);

  // Update video volume
  const video = document.querySelector("video");
  if (video) {
    video.volume = value;
  }

  // Update UI
  document.querySelector(".volume-control .value").textContent = value * 100 + "%";
  document.querySelector(".volume-control").classList.remove("open");
}

// Close menu when clicking outside
document.addEventListener("click", (e) => {
  if (!e.target.closest(".volume-control")) {
    document.querySelector(".volume-control").classList.remove("open");
  }
  if (!e.target.closest(".subtitle-control")) {
    document.querySelector(".subtitle-control").classList.remove("open");
  }
});

// Subtitle controls
function toggleSubtitleMenu(event) {
  if (event) {
    event.stopPropagation();
  }
  const control = document.querySelector(".subtitle-control");
  control.classList.toggle("open");
}

function setSubtitle(lang) {
  // Store in localStorage
  localStorage.setItem("subtitle", lang);

  // Update video subtitle tracks
  const video = document.querySelector("video");
  if (video && video.textTracks) {
    let langFound = false;
    for (let i = 0; i < video.textTracks.length; i++) {
      const track = video.textTracks[i];
      if (lang === "off") {
        track.mode = "disabled";
      } else if (track.language === lang) {
        track.mode = "showing";
        langFound = true;
      } else {
        track.mode = "disabled";
      }
    }

    // If requested language not found, default to 'off'
    if (lang !== "off" && !langFound) {
      localStorage.setItem("subtitle", "off");
      lang = "off";
    }
  }

  // Update UI only if control is not disabled
  const control = document.querySelector(".subtitle-control");
  const valueElement = document.querySelector(".subtitle-control .value");
  if (valueElement && control && !control.classList.contains("disabled")) {
    valueElement.textContent = lang === "off" ? "OFF" : lang.toUpperCase();
  }

  // Close menu
  document.querySelector(".subtitle-control").classList.remove("open");
}

function updateSubtitleButtonFromTracks(video) {
  if (!video || !video.textTracks) return;

  // Find the currently active subtitle track
  let activeLang = "off";
  for (let i = 0; i < video.textTracks.length; i++) {
    const track = video.textTracks[i];
    if (track.mode === "showing") {
      activeLang = track.language;
      break;
    }
  }

  // Update the button display
  const valueElement = document.querySelector(".subtitle-control .value");
  const control = document.querySelector(".subtitle-control");
  if (valueElement && control && !control.classList.contains("disabled")) {
    valueElement.textContent = activeLang === "off" ? "OFF" : activeLang.toUpperCase();
  }

  // Update localStorage if different from current setting
  const storedSubtitle = localStorage.getItem("subtitle") || "off";
  if (storedSubtitle !== activeLang) {
    localStorage.setItem("subtitle", activeLang);
  }
}

function populateSubtitleMenu(subtitles) {
  const control = document.querySelector(".subtitle-control");
  const button = document.querySelector("#subtitleButton");
  const menu = document.querySelector(".subtitle-menu");
  const valueElement = document.querySelector(".subtitle-control .value");

  if (!control || !button || !menu || !valueElement) return;

  // Clear existing options
  menu.innerHTML = "";

  if (!subtitles || subtitles.length === 0) {
    // No subtitles available - disable control
    control.classList.add("disabled");
    valueElement.textContent = "NA";
    button.disabled = true;
    button.onclick = null;
  } else {
    // Subtitles available - enable control
    control.classList.remove("disabled");
    button.disabled = false;
    button.onclick = toggleSubtitleMenu;

    // Add "Off" option
    const offButton = document.createElement("button");
    offButton.textContent = "Off";
    offButton.onclick = () => setSubtitle("off");
    menu.appendChild(offButton);

    // Add subtitle language options
    subtitles.forEach((sub) => {
      const button = document.createElement("button");
      button.textContent = sub.language.toUpperCase();
      button.onclick = () => setSubtitle(sub.language);
      menu.appendChild(button);
    });

    // Set current subtitle selection - check if stored preference is available
    const storedSubtitle = localStorage.getItem("subtitle") || "off";
    const availableLanguages = subtitles.map((sub) => sub.language);
    const currentLang =
      storedSubtitle !== "off" && availableLanguages.includes(storedSubtitle)
        ? storedSubtitle
        : "off";
    valueElement.textContent = currentLang === "off" ? "OFF" : currentLang.toUpperCase();

    // Update stored preference if it was invalid for this video
    if (currentLang !== storedSubtitle) {
      localStorage.setItem("subtitle", currentLang);
    }
  }
}

// Initialize subtitle on page load
document.addEventListener("DOMContentLoaded", () => {
  // Get stored volume or default to 1 (100%)
  const storedVolume = localStorage.getItem("volume");
  const volume = storedVolume !== null ? parseFloat(storedVolume) : 1;

  // Set initial volume display
  document.querySelector(".volume-control .value").textContent = volume * 100 + "%";

  // Set video volume if it exists
  const video = document.querySelector("video");
  if (video) {
    video.volume = volume;
  }
});

// Navigation
function toggleNavExpand() {
  const body = document.body;
  console.log("Current state:", {
    hidden: body.classList.contains("nav-hidden"),
    expanded: body.classList.contains("nav-expanded"),
  });

  if (body.classList.contains("nav-hidden")) {
    // From hidden (0%) to normal (100%)
    body.classList.remove("nav-hidden");
  } else if (body.classList.contains("nav-expanded")) {
    // From expanded (200%) to normal (100%)
    body.classList.remove("nav-expanded");
  } else {
    // From normal (100%) to expanded (200%)
    body.classList.add("nav-expanded");
  }

  console.log("New state:", {
    hidden: body.classList.contains("nav-hidden"),
    expanded: body.classList.contains("nav-expanded"),
  });
}

function closeNav() {
  const body = document.body;
  body.classList.add("nav-hidden");
  body.classList.remove("nav-expanded");
}

// Theme management
function toggleTheme() {
  const btn = document.getElementById("toggleTheme");
  const icon = btn.querySelector(".material-symbols-outlined");
  const isCurrentlyDark = !document.body.classList.contains("light-theme");

  if (isCurrentlyDark) {
    // Switch to light theme
    document.body.classList.add("light-theme");
    icon.textContent = "dark_mode";
    localStorage.setItem("theme", "light");
  } else {
    // Switch to dark theme
    document.body.classList.remove("light-theme");
    icon.textContent = "light_mode";
    localStorage.setItem("theme", "dark");
  }
}

// Autoplay functionality
let autoplayEnabled = localStorage.getItem("autoplay") === "true";

function toggleAutoplay() {
  autoplayEnabled = !autoplayEnabled;
  localStorage.setItem("autoplay", autoplayEnabled);
  updateAutoplayButton();
}

function updateAutoplayButton() {
  const btn = document.getElementById("toggleAutoplay");
  if (autoplayEnabled) {
    btn.classList.add("active");
    btn.title = "Autoplay ON";
  } else {
    btn.classList.remove("active");
    btn.title = "Autoplay OFF";
  }
}

function getNextVideo(currentPath) {
  const items = Array.from(document.querySelectorAll(".file-list li"));
  const currentIndex = items.findIndex((item) =>
    item.getAttribute("onclick")?.includes(currentPath),
  );
  if (currentIndex === -1) return null;

  // Look for next video (skip directories)
  for (let i = currentIndex + 1; i < items.length; i++) {
    const item = items[i];
    if (!item.querySelector(".duration")) continue; // Skip directories
    const onclick = item.getAttribute("onclick");
    if (!onclick) continue;
    const match = onclick.match(/playVideo\('([^']+)'\)/);
    if (match) return match[1];
  }
  return null;
}

// Video playback utilities
function updatePlayingClass(path) {
  // Remove playing class from previous video
  const previous = document.querySelector(".file-list li.playing");
  if (previous) {
    previous.classList.remove("playing");
  }

  // Add playing class to current video
  const current = document.querySelector(`.file-list li[onclick*="${CSS.escape(path)}"]`);
  if (current) {
    current.classList.add("playing");
  }
}

function createVideoElement(data, path) {
  let html =
    '<video controls autoplay><source src="/api/video/stream?path=' + path + '" type="video/mp4">';

  // Add subtitle tracks if available
  if (data.info.subtitles && data.info.subtitles.length > 0) {
    data.info.subtitles.forEach((sub) => {
      html +=
        `<track label="${sub.language}" kind="subtitles" srclang="${sub.language}" ` +
        `src="/api/video/subtitle?path=${path}&lang=${sub.language}">`;
    });
  }

  html += "</video>";
  return html;
}

function setupVideoElement(video, path) {
  // Set up all event listeners
  initVideoEventListeners(video);

  // Set stored volume
  const storedVolume = localStorage.getItem("volume");
  if (storedVolume !== null) {
    video.volume = parseFloat(storedVolume);
  }

  // Set stored subtitle
  const storedSubtitle = localStorage.getItem("subtitle") || "off";
  if (storedSubtitle !== "off") {
    // Wait for tracks to load, then apply subtitle setting
    video.addEventListener("loadedmetadata", () => {
      setTimeout(() => setSubtitle(storedSubtitle), 100);
    });
  }

  // Listen for text track changes (when user changes subtitle via browser controls)
  video.textTracks.addEventListener("change", () => {
    updateSubtitleButtonFromTracks(video);
  });

  // Also listen for when tracks are added
  video.addEventListener("addtrack", (event) => {
    if (event.track.kind === "subtitles") {
      updateSubtitleButtonFromTracks(video);
    }
  });

  // Add autoplay handling
  video.addEventListener("ended", () => {
    if (autoplayEnabled) {
      const nextPath = getNextVideo(path);
      if (nextPath) {
        playVideo(nextPath);
      }
    }
  });
}

function playVideo(path) {
  // Start video playback with subtitle info
  const player = document.getElementById("player");

  // Fetch video info first
  fetch("/api/video?path=" + path)
    .then((response) => response.json())
    .then((data) => {
      // Update playing class in file list
      updatePlayingClass(path);

      // Create and insert video element
      const videoHtml = createVideoElement(data, path);
      player.innerHTML = videoHtml;

      // Populate subtitle menu
      populateSubtitleMenu(data.info.subtitles);

      // Setup video element
      const video = player.querySelector("video");
      if (video) {
        setupVideoElement(video, path);
      }
    });
}

// Initialize global video event listeners
function initVideoEventListeners(video) {
  if (!video) return;

  // Sync play/pause button with video state
  video.addEventListener("play", () => updatePlayPauseButton(false));
  video.addEventListener("pause", () => updatePlayPauseButton(true));
  video.addEventListener("ended", () => updatePlayPauseButton(true));

  // Initial state
  updatePlayPauseButton(video.paused);
}

// Initialize on page load
document.addEventListener("DOMContentLoaded", () => {
  // Set theme (default to light theme)
  const theme = localStorage.getItem("theme") || "light";
  if (theme === "light") {
    document.body.classList.add("light-theme");
    document.getElementById("toggleTheme").querySelector(".material-symbols-outlined").textContent = "dark_mode";
  }

  // Initialize autoplay button state
  updateAutoplayButton();

  // Initialize subtitle control (disabled by default until video loads)
  const control = document.querySelector(".subtitle-control");
  const valueElement = document.querySelector(".subtitle-control .value");
  if (control && valueElement) {
    control.classList.add("disabled");
    valueElement.textContent = "NA";
  }

  // Initialize video controls if video exists
  initVideoEventListeners(document.querySelector("video"));
});
