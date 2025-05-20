import { getServerSession } from 'next-auth';
import { authConfig } from '../authConfig';
import {
  HistoryOffer,
  SaleOffer,
  SearchParams,
  UserProfile,
} from '../definitions';

export async function fetchSessionData(): Promise<UserProfile> {
  const user = (await getServerSession(authConfig))?.user;

  if (!user) {
    throw new Error('No session found');
  }

  const userProfile: UserProfile = {
    id: user.id,
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

  return userProfile;
}

export async function getHistory(
  params: SearchParams = {},
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
  params: SearchParams = {},
): Promise<HistoryOffer[]> {
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
  params: SearchParams = {},
): Promise<SaleOffer[]> {
  const url = `${process.env.URL}/api/account/favorites`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error('Failed to fetch favorite offers');
  }

  return response.json();
}

export async function fetchHistory(params: SearchParams) : Promise<{totalPages: number, totalOffers: number, offers: HistoryOffer[]}> {
    try{
        console.log("History page params: ", params);

        const data = await getHistory(params);

      
        console.log("History page data: ", data);

        return {
            totalPages: data.pagination.total_pages,
            totalOffers: data.pagination.total_records,
            offers: data.offers
        };


    } catch (error) {
        console.error("Api error:", error);
        throw new Error('Failed to fetch home page data.');
    }
}
