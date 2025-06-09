import { getServerSession } from 'next-auth';
import { cache } from 'react';
import { fetchWithRefresh } from '../../api/fetchWithRefresh';
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
    filter: {},
  },
): Promise<{
  offers: HistoryOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/purchased-offers`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch user's listings");
  }

  return response.json();
}

export async function getListings(params: SearchParams): Promise<{
  offers: SaleOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/my-offers`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    },
  );

  if (!response.ok) {
    throw new Error("Failed to fetch user's listings");
  }

  return response.json();
}

export async function getFavorites(
  params: SearchParams = {
    pagination: {
      page: 1,
      page_size: 6,
    },
    filter: {},
  },
): Promise<{
  offers: SaleOffer[];
  pagination: {
    total_pages: number;
    total_records: number;
  };
}> {
  const response = await fetchWithRefresh(
    `${process.env.API_URL}/sale-offer/liked-offers`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(params),
    },
  );

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

    const mappedResponse: {
      totalPages: number;
      totalOffers: number;
      offers: HistoryOffer[];
    } = {
      offers: data.offers.map((offer: any) => ({
        id: offer.id,
        name: offer.name,
        production_year: offer.production_year,
        mileage: offer.mileage,
        color: offer.color,
        price: offer.price,
        is_auction: offer.is_auction,
        main_url: offer.main_url,
        can_modify: offer.can_modify,
        date_end: offer.issue_date,
        seller_id: offer.seller_id,
        seller_name: offer.username,
      })),
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
    };

    return mappedResponse;
  } catch (error) {
    console.error('Api error:', error);
    throw new Error("Failed to fetch user's history.");
  }
}

export async function fetchListings(params: SearchParams): Promise<{
  totalPages: number;
  totalOffers: number;
  offers: SaleOffer[];
}> {
  try {
    const data = await getListings(params);

    return {
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
      offers: data.offers,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error("Failed to fetch user's listings.");
  }
}

export async function fetchFavorites(
  params: SearchParams,
): Promise<{ totalPages: number; totalOffers: number; offers: SaleOffer[] }> {
  try {
    const data = await getFavorites(params);

    return {
      totalPages: data.pagination.total_pages,
      totalOffers: data.pagination.total_records,
      offers: data.offers,
    };
  } catch (error) {
    console.error('Api error:', error);
    throw new Error('Failed to fetch favorite offers.');
  }
}
