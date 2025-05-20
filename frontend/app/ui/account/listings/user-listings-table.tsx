import { getListings } from '@/app/lib/account/data';
import { BaseOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffersTable from '@/app/ui/generic-offer-table';
import SingleListingsOffer from './listings-offer-card';

export default async function UsersListingsTable({
  params,
}: {
  params: SearchParams;
}) {
  return (
    <GenericOffersTable<BaseOffer>
      fetchFunction={getListings}
      params={params}
      ItemComponent={SingleListingsOffer}
    />
  );
}
