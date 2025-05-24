export default async function Page(props: { params: Promise<{id: string }> }) {
    const { params } = props;
    const { id } = await params;

    return (
        <p className="text-2xl font-bold text-center">
            Offer ID: {id}
        </p>
    )
}