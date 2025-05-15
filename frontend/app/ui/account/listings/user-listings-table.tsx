import { fetchListings } from '@/app/lib/account/data';
import { SaleOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffersTable from '@/app/ui/generic-offer-table';
import SingleListingsOffer from './single-listings-offer';

export default async function UsersListings({
  params,
}: {
  params: SearchParams;
}) {
  return (
    <GenericOffersTable<SaleOffer>
      fetchFunction={fetchListings}
      params={params}
      ItemComponent={SingleListingsOffer}
    />
  );
}
