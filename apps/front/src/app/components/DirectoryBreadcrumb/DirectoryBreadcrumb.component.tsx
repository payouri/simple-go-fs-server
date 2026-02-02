import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
} from '@/components/ui/breadcrumb';
import { ChevronRight } from 'lucide-react';
import { Link } from 'react-router-dom';

export type BreadcrumbsProps = {
  path: string[];
  separator?: string;
  showHome?:
    | {
        text: string;
        href: string;
      }
    | false;
  // onLinkClick?: (path: string[]) => void;
};

const defaultSeparator = <ChevronRight />;

export function DirectoryBreadcrumb(props: BreadcrumbsProps) {
  const {
    path,
    separator = defaultSeparator,
    showHome = {
      text: 'Home',
      href: '/',
    },
  } = props;

  console.log({ path });

  return (
    <Breadcrumb>
      <BreadcrumbList>
        {showHome ? (
          <BreadcrumbItem>
            <BreadcrumbLink href="#">
              <Link to={showHome.href} relative="route">
                <BreadcrumbPage>{showHome.text}</BreadcrumbPage>
              </Link>
            </BreadcrumbLink>
          </BreadcrumbItem>
        ) : null}
        {path.length > 0 && showHome ? (
          <BreadcrumbSeparator>{separator}</BreadcrumbSeparator>
        ) : null}
        {path.map((item, index) => (
          <>
            <BreadcrumbItem key={`${item}-${index}`}>
              <BreadcrumbLink href="#">
                <Link
                  to={{
                    pathname: path.slice(0, index + 1).join('/'),
                  }}
                  relative="route"
                >
                  <BreadcrumbPage>{item}</BreadcrumbPage>
                </Link>
              </BreadcrumbLink>
            </BreadcrumbItem>
            {index < path.length - 1 && (
              <BreadcrumbSeparator>{separator}</BreadcrumbSeparator>
            )}
          </>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
