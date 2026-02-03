import {
  Breadcrumb,
  BreadcrumbList,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbSeparator,
  BreadcrumbPage,
  BreadcrumbReactRouterLink,
} from '@/components/ui/breadcrumb';
import { ChevronRight } from 'lucide-react';
import { Link } from 'react-router-dom';
import { Fragment } from 'react/jsx-runtime';

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

  return (
    <Breadcrumb className="flex-[0_0_auto]">
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
          <Fragment key={`${item}-${index}`}>
            <BreadcrumbItem>
              <BreadcrumbReactRouterLink
                to={{
                  pathname: path.slice(0, index + 1).join('/'),
                }}
                relative="route"
              >
                <BreadcrumbPage>{item}</BreadcrumbPage>
              </BreadcrumbReactRouterLink>
            </BreadcrumbItem>
            {index < path.length - 1 && (
              <BreadcrumbSeparator>{separator}</BreadcrumbSeparator>
            )}
          </Fragment>
        ))}
      </BreadcrumbList>
    </Breadcrumb>
  );
}
