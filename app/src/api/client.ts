export async function apiFetch(
  url: string,
  token: string
) {

  return fetch(url, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
}
