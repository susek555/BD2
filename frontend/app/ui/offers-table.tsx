import { SaleOffer, SearchParams } from '@/app/lib/definitions';
import GenericOffersTable from './generic-offer-table';
import SaleOfferCard from './sale-offer-card';

export default async function OffersTable({
  params,
  fetchFunction,
}: {
  params: SearchParams;
  fetchFunction: (params: SearchParams) => Promise<SaleOffer[]>;
}) {
  return (
    <GenericOffersTable<SaleOffer>
      fetchFunction={fetchFunction}
      params={params}
      ItemComponent={SaleOfferCard}
    />
  );
}
