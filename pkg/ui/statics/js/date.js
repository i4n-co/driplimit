function addDays(date, delta) {
  return addHours(date, 24 * delta);
}
function addHours(date, delta) {
  return addMinutes(date, 60 * delta);
}
function addMinutes(date, delta) {
  return addSeconds(date, 60 * delta);
}
function addSeconds(date, delta) {
  return new Date(date.getTime() + 1000 * delta);
}

function datetimeFormat(date) {
  function ensureLeadingZero(n) {
    return `${n}`.padStart(2, 0);
  }
  const year = date.getFullYear();
  const month = ensureLeadingZero(date.getMonth() + 1);
  const day = ensureLeadingZero(date.getDate());
  const hours = ensureLeadingZero(date.getHours());
  const minutes = ensureLeadingZero(date.getMinutes());
  return `${year}-${month}-${day}T${hours}:${minutes}`;
}

function localtz() {
  return Intl.DateTimeFormat().resolvedOptions().timeZone;
}

function timer(duration) {
  return {
    start: Date.now(),
    time: Date.now(),
    closed: false,
    duration: duration,
    init() {
      setInterval(() => {
        this.time = Date.now();
      }, 100);
      return this;
    },
    done() {
      return closed || this.time >= this.start + this.duration;
    },
    restart(duration) {
      this.start = Date.now();
      this.duration = duration
    },
    close() {
      this.closed = true;
    },
  };
}
