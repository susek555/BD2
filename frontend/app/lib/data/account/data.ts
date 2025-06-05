import { getServerSession } from 'next-auth';
import { cache } from 'react';
import { authConfig } from '../../authConfig';
import { HistoryOffer, SaleOffer } from '../../definitions/SaleOffer';
import { SearchParams } from '../../definitions/SearchParams';
import { UserProfile } from '../../definitions/user';

export async function fetchSessionData(): Promise<UserProfile> {
  const user = (await getServerSession(authConfig))?.user;

  if (!user) {
    throw new Error('No session found');
  }

  const userProfile: UserProfile = {
    selector: user.selector,
    userId: user.userId,
    username: user.username,
    email: user.email,
  };

  if (user.selector === 'P') {
    userProfile.personName = user.personName;
    userProfile.personSurname = user.personSurname;
  } else if (user.selector === 'C') {
    userProfile.companyName = user.companyName;
    userProfile.companyNip = user.companyNip;
  }
  console.log(userProfile);

  return userProfile;
}

export const cachedSessionData = cache(async () => {
  return fetchSessionData();
});

export async function getHistory(
  params: SearchParams = {
    pagination: {
      page: 1,
      page_size: 6,
    },
  },
): Promise<{
  offers: HistoryOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const url = `${process.env.URL}/api/account/history`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  return response.json();
}

export async function getListings(
  params: SearchParams = {
    pagination: {
      page: 1,
      page_size: 6,
    },
  },
): Promise<{
  offers: HistoryOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const url = `${process.env.URL}/api/account/listings`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch history offers');
  }

  return response.json();
}

export async function getFavorites(
  params: SearchParams = {
    pagination: {
      page: 1,
      page_size: 6,
    },
  },
): Promise<{
  offers: SaleOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const url = `${process.env.URL}/api/account/favorites`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch favorite offers');
  }

  return response.json();
}

export async function fetchHistory(params: SearchParams): Promise<{
  totalPages: number;
  totalOffers: number;
  offers: HistoryOffer[];
}> {
  try {
    const data = await getHistory(params);

    console.log('History page data: ', data);

    return {
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
      offers: data.offers,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error('Failed to fetch home page data.');
  }
}

export async function fetchListings(params: SearchParams): Promise<{
  totalPages: number;
  totalOffers: number;
  offers: HistoryOffer[];
}> {
  try {
    const data = await getListings(params);

    console.log('Listings page data: ', data);

    return {
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
      offers: data.offers,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error('Failed to fetch home page data.');
  }
}

export async function fetchFavorites(
  params: SearchParams,
): Promise<{ totalPages: number; totalOffers: number; offers: SaleOffer[] }> {
  try {
    const data = await getFavorites(params);

    console.log('Favorites page data: ', data);

    return {
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
      offers: data.offers,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error('Failed to fetch home page data.');
  }
}
