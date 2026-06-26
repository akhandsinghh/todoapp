export function unwrapApiResponse(payload) {
  if (payload === null || payload === undefined) return payload;
  if (Array.isArray(payload)) return payload;
  if (typeof payload !== 'object') return payload;

  if (Object.prototype.hasOwnProperty.call(payload, 'model')) {
    return unwrapApiResponse(payload.model);
  }

  if (Object.prototype.hasOwnProperty.call(payload, 'data')) {
    return unwrapApiResponse(payload.data);
  }

  return payload;
}

export function normalizeListResponse(payload) {
  const unwrapped = unwrapApiResponse(payload);

  if (Array.isArray(unwrapped)) {
    return { items: unwrapped, total: unwrapped.length };
  }

  if (unwrapped && typeof unwrapped === 'object') {
    const items = Array.isArray(unwrapped.items)
      ? unwrapped.items
      : Array.isArray(unwrapped.data)
        ? unwrapped.data
        : [];

    const total = Number.isFinite(unwrapped.total)
      ? unwrapped.total
      : Number.isFinite(unwrapped.count)
        ? unwrapped.count
        : items.length;

    return { items, total };
  }

  return { items: [], total: 0 };
}

export function normalizeArrayResponse(payload) {
  return normalizeListResponse(payload).items;
}
