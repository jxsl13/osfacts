# osfacts

Small library that allows to detect system information like operating system, its version and derivate



## Reference files


### Bash Reference

```shell
# Determine OS, product, revision
   str_OS=$(uname -s)
   if [ "$str_OS" = "AIX" ]; then
         str_DIST="$str_OS"
         str_REV=" $(oslevel | cut -c1-3)"
   elif [ -x /SZIR/bin/osinfo ]; then
      str_OSINFO="/SZIR/bin/osinfo";
      str_PRODUCT="$($str_OSINFO -o)"
      if [ "$str_PRODUCT" = "SLES" ]; then
         str_DIST=$($str_OSINFO -om | tr -d "\n")
         str_REV=" SP$($str_OSINFO -p | tr -d "\n")"
      elif [ "$str_PRODUCT" = "CENTOS" ]; then
         str_DIST=$($str_OSINFO -o | tr -d "\n")
         str_REV=" $($str_OSINFO -m | tr -d "\n").$($str_OSINFO -p | tr -d "\n")"
      elif [ "$str_PRODUCT" = "RHEL" ] || [ "$str_PRODUCT" = "ORACLE" ] || [ "$str_PRODUCT" = "ORACLE EXADATA" ]; then
         str_DIST=$($str_OSINFO -o | tr -d "\n")
         str_REV=" $($str_OSINFO -m | tr -d "\n").$($str_OSINFO -p | tr -d "\n")"
      elif [ "$str_PRODUCT" = "SLES4SAP" ]; then
         str_DIST=$($str_OSINFO -om | tr -d "\n")
         str_REV=" SP$($str_OSINFO -p | tr -d "\n")"
      elif [ "$str_PRODUCT" = "SunOS" ]; then
         str_DIST="Solaris ";
         str_REV="$($str_OSINFO -p | tr -d "\n")"
         if [ "$str_REV" == "11" ]; then
            str_REV="$(uname -v)";
         fi
      elif [ "$str_OS" = "HP-UX" ]; then
         str_DIST="$(uname)"
         str_REV=" $(uname -r)"
      fi
   elif [ -f /etc/oracle-release ]; then
      str_PRODUCT="ORACLE"
      str_DIST="ORACLE "
      str_REV=" $(cat /etc/oracle-release | awk '{ print $5 }')"
   fi
```