import { fetchOffers } from '@/app/lib/data';
import { SaleOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffer from '../generic-offer';
import GenericOffersTable from '../generic-offer-table';

export default async function OffersTable({
  params,
}: {
  params: SearchParams;
}) {
  return (
    <GenericOffersTable<SaleOffer>
      fetchFunction={fetchOffers}
      params={params}
      ItemComponent={GenericOffer}
    />
  );
}
