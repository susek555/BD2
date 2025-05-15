import { fetchOffers } from '@/app/lib/data';
import { SaleOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffersTable from '../generic-offer-table';
import SingleOffer from './single-offer';

export default async function OffersTable({
  params,
}: {
  params: SearchParams;
}) {
  return (
    <GenericOffersTable<SaleOffer>
      fetchFunction={fetchOffers}
      params={params}
      ItemComponent={SingleOffer}
    />
  );
}
