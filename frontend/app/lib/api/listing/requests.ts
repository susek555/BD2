import { getApiUrl } from '../../get-api-url';

export const fetchListing = async (id: string) => {
  const response = await fetch(`/api/listing/${id}`, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch listing');
  }

  return response.json();
};

// TODO create listing type with all data from backend and update listingData type
export const updateListing = async (id: string, listingData: any) => {
  const response = await fetch(`/api/listing/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(listingData),
  });

  if (!response.ok) {
    throw new Error('Failed to update listing');
  }

  return response.json();
};

export const createListing = async (listingData: any) => {
  const response = await fetch('/api/listing/create', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(listingData),
  });

  if (!response.ok) {
    throw new Error('Failed to create listing');
  }

  return response.json();
};

export const deleteListing = async (id: string) => {
  const url = getApiUrl(`/api/listing/${id}`);
  console.log(url);

  const response = await fetch(url, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error('Failed to delete listing');
  }

  return response.json();
};

export const updateFavoriteStatus = async (id: string, isFavorite: boolean) => {
  let response;
  if (isFavorite) {
    response = await fetch(`/api/favorites/${id}`, {
      method: 'DELETE',
    });
  } else {
    response = await fetch(`/api/favorites/${id}`, {
      method: 'POST',
    });
  }

  if (!response.ok) {
    throw new Error('Failed to add listing to favorites');
  }

  return response.json();
};
