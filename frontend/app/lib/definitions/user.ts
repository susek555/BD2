// Account page definitions

export interface Tab {
  id: string;
  label: string;
  href: string;
}

export const profileTabs: Tab[] = [
  { id: 'activity', label: 'Activity', href: '/account/activity' },
  { id: 'listings', label: 'Listings', href: '/account/listings' },
  { id: 'favorites', label: 'Favorites', href: '/account/favorites' },
  { id: 'reviews', label: 'Reviews', href: '/account/reviews' },
  { id: 'settings', label: 'Settings', href: '/account/settings' },
];

export interface UserProfile {
  id: number;
  username: string;
  email: string;
  personName?: string;
  personSurname?: string;
  companyName?: string;
  companyNip?: string;
}