%global debug_package %{nil}

Name: @NAME@
Summary: @SUMMARY@
Version: @VERSION@
Release: 1%{?dist}
License: MIT

Source0: %{name}
Source1: %{name}.desktop
Source2: %{name}.png

%description
@DESCRIPTION@

%install
mkdir -p %{buildroot}%{_bindir} %{buildroot}%{_datadir}/{applications,pixmaps}

cp %{SOURCE0} %{buildroot}%{_bindir}/%{name}
cp %{SOURCE1} %{buildroot}%{_datadir}/applications/%{name}.desktop
cp %{SOURCE2} %{buildroot}%{_datadir}/pixmaps/%{name}.png

%post
# Install the desktop entry
update-desktop-database &> /dev/null || :

%postun
# Uninstall the desktop entry
update-desktop-database &> /dev/null || :


%files
%{_bindir}/%{name}
%{_datadir}/applications/%{name}.desktop
%{_datadir}/pixmaps/%{name}.png


%changelog
* Sat Feb 17 2024 @AUTHOR@ <@AUTHOR_EMAIL@> - %{version}-1
- Initial build